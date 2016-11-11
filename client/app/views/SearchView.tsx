import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { Grid, Section, Row, Col } from '../components/Grid';
import { Pagination, PageButton } from '../components/Pagination';
import { ExtFacet } from '../components/ExtFacet';
import { FacetPanel } from '../components/FacetPanel';
import { FullRefsFacet } from '../components/FullRefsFacet';
import { Facets } from '../components/Facets';
import { FileContent } from '../components/FileContent';
import { RootState, SearchResult, SearchFacets, FilterParams } from '../reducers';
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
    handleExtToggle = (term: string) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, mergeTerm('x', this.props.filterParams, term));
    };

    handleOrganizationToggle = (term: string) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, mergeTerm('o', this.props.filterParams, term));
    };

    handleProjectToggle = (term: string) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, mergeTerm('p', this.props.filterParams, term));
    };

    handleRepositoryToggle = (term: string) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, mergeTerm('r', this.props.filterParams, term));
    };

    handleBranchesToggle = (term: string) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, mergeTerm('b', this.props.filterParams, term));
    };

    handleTagsToggle = (term: string) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, mergeTerm('t', this.props.filterParams, term));
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
        const { loading, filterParams, result, facets } = this.props;

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
                        <FacetPanel title='File extensions'
                            facet={facets.facets['ext']}
                            emptyKeyword='/noext/'
                            emptyLabel='(No extension)'
                            selected={filterParams.x}
                            onToggle={this.handleExtToggle} />
                        <FacetPanel title='Organization'
                            facet={facets.facets['organization']}
                            selected={filterParams.o}
                            onToggle={this.handleOrganizationToggle} />
                        <FacetPanel title='Project'
                            facet={facets.facets['project']}
                            selected={filterParams.p}
                            onToggle={this.handleProjectToggle} />
                        <FacetPanel title='Repository'
                            facet={facets.facets['repository']}
                            selected={filterParams.r}
                            onToggle={this.handleRepositoryToggle} />
                        <FacetPanel title='Branches'
                            facet={facets.facets['branches']}
                            selected={filterParams.b}
                            onToggle={this.handleBranchesToggle} />
                        <FacetPanel title='Tags'
                            facet={facets.facets['tags']}
                            selected={filterParams.t}
                            onToggle={this.handleTagsToggle} />

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

function mergeTerm(key: string, params: FilterParams, term: string) {
    const target = params[key] || [];
    let terms;
    if (target.find(x => x === term)) {
        terms = target.filter(x => x !== term);
    } else {
        terms = target.concat(term);
    }
    return Object.assign({}, params, {
        [key]: terms
    });
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