import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { Grid, Section, Row, Col } from '../components/Grid';
import { SearchSidePanel } from '../components/SearchSidePanel';
import { Pagination, PageButton } from '../components/Pagination';
import { ExtFacet } from '../components/ExtFacet';
import { FacetPanel } from '../components/FacetPanel';
import { FullRefsFacet } from '../components/FullRefsFacet';
import { Facets } from '../components/Facets';
import { FileContent } from '../components/FileContent';
import { RootState, SearchResult, SearchFacets, FilterParams, FacetKey } from '../reducers';
import * as Actions from '../actions';

const MDSpinner = require('react-md-spinner').default;

interface Props {
    dispatch?: Dispatch<Action>;
    loading: boolean;
    query: string;
    filterParams: FilterParams;
    result: SearchResult;
    facets: SearchFacets;
}

class SearchView extends React.Component<Props, void> {
    handleFacetToggle = (filterParams: FilterParams) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, filterParams);
    };

    showPage = (page: number) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, this.props.filterParams, page);
    };

    next = () => {
        let next = this.props.result.current + 1;
        const pageSize = Math.ceil(this.props.result.size / this.props.result.limit) || 0;
        if (next >= pageSize) {
            next = pageSize - 1;
        }
        Actions.searchFilter(this.props.dispatch, this.props.query, this.props.filterParams, next);
    };

    prev = () => {
        let next = this.props.result.current - 1;
        if (next < 0) {
            next = 0;
        }
        Actions.searchFilter(this.props.dispatch, this.props.query, this.props.filterParams, next);
    };

    render() {
        const { loading, filterParams, result, facets, query } = this.props;

        const pageSize = Math.ceil(result.size / result.limit) || 0;
        const pageButtons = [];

        let start = result.current - 2;
        let end = result.current + 2;
        const lastIndex = pageSize - 1;

        if (start < 0) {
            end += -start;
            start = 0;
        }
        if (end > lastIndex) {
            if (start > 0) {
                start -= (end - lastIndex);
                if (start < 0) {
                    start = 0;
                }
            }
            end = lastIndex;
        }

        if (start > 0) {
            pageButtons.push((
                <li>
                    <PageButton onClick={this.showPage.bind(null, 0)}>1</PageButton>
                </li>
            ));
            pageButtons.push(<li>...</li>);
        }

        for (let index = start; index <= end; index++) {
            pageButtons.push((
                <li>
                    <PageButton onClick={this.showPage.bind(null, index)} isActive={index === result.current}>{index + 1}</PageButton>
                </li>
            ));
        }

        if (end < lastIndex) {
            pageButtons.push(<li>...</li>);
            pageButtons.push((
                <li>
                    <PageButton onClick={this.showPage.bind(null, lastIndex)}>{lastIndex + 1}</PageButton>
                </li>
            ));
        }

        return (
            <Section>
                <Row>
                    <Col size='is3'>
                        <SearchSidePanel facets={facets}
                            filterParams={filterParams}
                            query={query}
                            onToggle={this.handleFacetToggle} />
                    </Col>
                    <Col size='is9'>
                        <Grid>
                            {result.hits.map(x => {
                                return (
                                    <Row key={x.blob}>
                                        <Col size='is12'>
                                            <FileContent metadata={x} keyword={x.keyword} preview={x.preview} />
                                        </Col>
                                    </Row>
                                );
                            })}
                            {result && result.size > 10 &&
                                <Row>
                                    <Col size='is12'>
                                        <Pagination>
                                            <PageButton onClick={this.prev}>Previous</PageButton>
                                            <PageButton onClick={this.next}>Next page</PageButton>
                                            <ul>
                                                {pageButtons}
                                            </ul>
                                        </Pagination>
                                    </Col>
                                </Row>
                            }
                        </Grid>
                    </Col>
                </Row>
            </Section>
        );
    }
}

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        loading: state.app.present.loading,
        query: state.app.present.query,
        filterParams: state.app.present.filterParams,
        result: state.app.present.result,
        facets: state.app.present.facets
    };
}

const SearchViewContainer = connect(
    mapStateToProps
)(SearchView);

export default SearchViewContainer;