import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { Grid, Section, Row, Col } from '../components/Grid';
import { SearchSidePanel } from '../components/SearchSidePanel';
import { SearchResultPanel } from '../components/SearchResultPanel';
import { ExtFacet } from '../components/ExtFacet';
import { FacetPanel } from '../components/FacetPanel';
import { FullRefsFacet } from '../components/FullRefsFacet';
import { Facets } from '../components/Facets';
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

    handlePageChange = (page: number) => {
        Actions.searchFilter(this.props.dispatch, this.props.query, this.props.filterParams, page);
    };

    render() {
        const { loading, filterParams, result, facets, query } = this.props;

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
                        <SearchResultPanel result={result} onPageChange={this.handlePageChange} />
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