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
import { RootState, SearchResult, SearchFacets, BaseFilterParams, FilterParams, FacetKey } from '../reducers';
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
    organizations?: string[];
    projects?: string[];
    repositories?: string[];
}

class SearchView extends React.Component<Props, void> {
    componentWillMount() {
        let count = 0;
        this.props.history.listen((arg1, {location, params}) => {
            if (location.query.q !== undefined && location.query.q !== '') {
                Actions.search(this.props.dispatch, params, location.query);
            }
        });
        Actions.getBaseFilters(this.props.dispatch,
            this.props.params.organization,
            this.props.params.project,
            this.props.params.repository);
    }

    handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.keyCode === 13) {
            Actions.triggerSearch(this.props.dispatch, e.target['value']);
        }
    };

    handleFacetToggle = (filterParams: FilterParams) => {
        Actions.triggerFilter(this.props.dispatch, filterParams);
    };

    handlePageChange = (page: number) => {
        Actions.triggerFilter(this.props.dispatch, this.props.filterParams, page);
    };

    render() {
        const { loading, query, filterParams, result, facets,
            organizations, projects, repositories } = this.props;

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
                    organizations={organizations}
                    projects={projects}
                    repositories={repositories}
                    organization={this.props.params.organization}
                    project={this.props.params.project}
                    repository={this.props.params.repository}
                    branch={this.props.params.branch}
                    tag={this.props.params.tag}
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
        organizations: state.app.present.baseFilterOptions.organizations,
        projects: state.app.present.baseFilterOptions.projects,
        repositories: state.app.present.baseFilterOptions.repositories,
        filterParams: props.location.query,
        result: state.app.present.result,
        facets: state.app.present.facets
    };
}

const SearchViewContainer = connect(
    mapStateToProps
)(SearchView);

export default SearchViewContainer;