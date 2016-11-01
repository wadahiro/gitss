import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { Grid, Row, Col } from '../components/Grid';
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
        Actions.search(this.props.dispatch, this.props.query, mergeTerm('x', this.props.filterParams, term));
    };

    handleOrganizationToggle = (term: string) => {
        Actions.search(this.props.dispatch, this.props.query, mergeTerm('o', this.props.filterParams, term));
    };

    handleProjectToggle = (term: string) => {
        Actions.search(this.props.dispatch, this.props.query, mergeTerm('p', this.props.filterParams, term));
    };

    handleRepositoryToggle = (term: string) => {
        Actions.search(this.props.dispatch, this.props.query, mergeTerm('r', this.props.filterParams, term));
    };

    handleRefsToggle = (term: string) => {
        Actions.search(this.props.dispatch, this.props.query, mergeTerm('b', this.props.filterParams, term));
    };

    render() {
        const { loading, filterParams, result, facets } = this.props;
        return (
            <Grid>
                <Row>
                    <Col size='is3'>
                        <FacetPanel title='File extensions'
                            facet={facets.facets['ext']}
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
                            facet={facets.facets['refs']}
                            selected={filterParams.b}
                            onToggle={this.handleRefsToggle} />

                    </Col>
                    <Col size='is9'>
                        <Grid>
                            {result.hits.map(x => {
                                return (
                                    <Row key={x._source.blob}>
                                        <Col size='is12'>
                                            <FileContent metadata={x._source.metadata} keyword={x.keyword} preview={x.preview} />
                                        </Col>
                                    </Row>
                                );
                            })}
                        </Grid>
                    </Col>
                </Row>
            </Grid>
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