import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { NavBar } from '../components/NavBar';
import { Grid, Section, Row, Col } from '../components/Grid';
import { SearchSidePanel } from '../components/SearchSidePanel';
import { SearchResultPanel } from '../components/SearchResultPanel';
import { Scrollbars } from '../components/Scrollbars';
import { ExtFacet } from '../components/ExtFacet';
import { FacetPanel } from '../components/FacetPanel';
import { FullRefsFacet } from '../components/FullRefsFacet';
import { Facets } from '../components/Facets';
import { RootState, SearchResult, SearchFacets, BaseFilterParams, BaseFilterOptions, FilterParams, FacetKey } from '../reducers';
import * as Actions from '../actions';

interface Props {
    query: string;
    loading: boolean;
    filterParams: FilterParams;
    result: SearchResult;
    facets: SearchFacets;
    // react-redux inject props
    dispatch?: Dispatch<Action>;
    // react-router inject props
    location?: any;
    history?: any;
    params?: BaseFilterParams;
    // lazy fetch
    baseFilterParams?: BaseFilterParams;
    baseFilterOptions?: BaseFilterOptions;
}

class SearchView extends React.Component<Props, void> {
    componentWillMount() {
        let count = 0;
        this.props.history.listen((arg1, {location, params}) => {
            if (location.query.q !== undefined && location.query.q !== '') {
                Actions.search(this.props.dispatch, params, location.query);
            }

            Actions.getBaseFilters(this.props.dispatch,
                params.organization,
                params.project,
                params.repository);
        });
    }

    handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.keyCode === 13) {
            Actions.triggerSearch(this.props.dispatch, this.props.baseFilterParams, e.target['value']);
        }
    };

    handleFacetToggle = (filterParams: FilterParams) => {
        Actions.triggerFilter(this.props.dispatch, this.props.baseFilterParams, filterParams, this.props.query);
    };

    handlePageChange = (page: number) => {
        Actions.triggerFilter(this.props.dispatch, this.props.baseFilterParams, this.props.filterParams, this.props.query, page);
    };

    handleBaseFilterChange = (values: BaseFilterParams) => {
        Actions.triggerBaseFilter(this.props.dispatch, values, this.props.filterParams, this.props.query);
    };

    render() {
        const { loading, query, filterParams, result, facets,
            baseFilterParams, baseFilterOptions } = this.props;

        const sidePanelStyle = {
            position: 'fixed',
            width: 300,
            hight: 700
        };
        const resultPanelStyle = {
            paddingLeft: 320,
            width: '100%'
        };

        return (
            <div>
                <NavBar onKeyDown={this.handleKeyDown}
                    loading={this.props.loading}
                    result={this.props.result}
                    query={query}
                    baseFilterParams={baseFilterParams}
                    baseFilterOptions={baseFilterOptions}
                    onBaseFilterChange={this.handleBaseFilterChange}
                    />
                <Section style={{ marginTop: 80 }}>
                    <Row>
                        <Col size='is3' style={sidePanelStyle}>
                            <Scrollbars style={{ height: 600 }}>
                                <SearchSidePanel facets={facets}
                                    searchParams={filterParams}
                                    onToggle={this.handleFacetToggle} />
                            </Scrollbars>
                        </Col>
                        <Col size='is9' style={resultPanelStyle}>
                            <SearchResultPanel result={result} onPageChange={this.handlePageChange} />
                        </Col>
                    </Row>
                </Section>
            </div>
        );
    }
}

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        loading: state.app.present.loading,
        query: props.location.query['q'] !== undefined ? props.location.query['q'] : '',
        filterParams: props.location.query,
        result: state.app.present.result,
        facets: state.app.present.facets,

        // Convert react-router injected params
        baseFilterParams: toBaseFilterParams(props.params),
        baseFilterOptions: state.app.present.baseFilterOptions
    };
}

function toBaseFilterParams(params: Object): BaseFilterParams {
    return params;
}

const SearchViewContainer = connect(
    mapStateToProps
)(SearchView);

export default SearchViewContainer;