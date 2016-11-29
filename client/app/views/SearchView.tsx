import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { NavBar } from '../components/NavBar';
import { AppFooter } from '../components/Footer';
import { Grid, Section, Row, Col, StickyFooterPage } from '../components/Grid';
import { SearchSidePanel } from '../components/SearchSidePanel';
import { SearchResultPanel } from '../components/SearchResultPanel';
import { SideBar } from '../components/SideBar';
import { ExtFacet } from '../components/ExtFacet';
import { FacetPanel } from '../components/FacetPanel';
import { FullRefsFacet } from '../components/FullRefsFacet';
import { Facets } from '../components/Facets';
import { RootState, SearchResult, SearchFacets, BaseFilterParams, BaseFilterOptions, FilterParams, FacetKey } from '../reducers';
import * as Actions from '../actions';

interface Props {
    query: string;
    loading: boolean;
    showOptions: boolean;
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
    unlisten = null;

    componentWillMount() {
        let count = 0;
        this.unlisten = this.props.history.listen((arg1, { location }) => {
            if (location.query.q !== undefined && location.query.q !== '') {
                Actions.search(this.props.dispatch, location.query);
            }
        });
    }

    componentWillUnmount() {
        this.unlisten();
    }

    handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.keyCode === 13) {
            Actions.triggerSearch(this.props.dispatch, e.target['value']);
        }
    };

    handleFacetToggle = (filterParams: FilterParams) => {
        Actions.triggerFilter(this.props.dispatch, filterParams, this.props.query);
    };

    handlePageChange = (page: number) => {
        Actions.triggerFilter(this.props.dispatch, this.props.filterParams, this.props.query, page);
    };

    handleSideBarToggle = () => {
        Actions.toggleSearchOptions(this.props.dispatch);
    };

    render() {
        const { loading, showOptions, query, filterParams, result, facets,
            baseFilterParams, baseFilterOptions } = this.props;

        return (
            <div>
                <NavBar onKeyDown={this.handleKeyDown}
                    showSideBarToggle={!showOptions}
                    onSideBarToggleClick={this.handleSideBarToggle}
                    loading={this.props.loading}
                    result={this.props.result}
                    query={query}
                    />
                <Row isGapless style={{ marginTop: 88 }}>
                    <Col size='isNarrow'>
                        <div style={{ width: 280 }}>
                            <SearchSidePanel style={sidePanelStyle}
                                facets={facets}
                                searchParams={filterParams}
                                onToggle={this.handleFacetToggle} />
                        </div>
                    </Col>
                    <Col>
                        <div style={{ paddingTop: 20 }}>
                            <StickyFooterPage footer={<AppFooter />}>
                                <SearchResultPanel style={resultPanelStyle}
                                    result={result}
                                    onPageChange={this.handlePageChange} />
                            </StickyFooterPage>
                        </div>
                    </Col>
                </Row>
            </div>
        );
    }
}

const sidePanelStyle = {
    position: 'fixed',
    height: 'calc(100vh - 108px)',
    overflowY: 'auto',
    backgroundColor: 'white',
    borderRight: '1px solid rgb(204, 204, 204)',
    // top: 100,
    // left: 0,
    width: 280,
    padding: 10
};

const resultPanelStyle = {
    flex: '1 1',
    paddingLeft: 50,
    paddingRight: 30,
    paddingBottom: 50
};

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        loading: state.app.present.loading,
        showOptions: state.app.present.showSearchOptions,
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