import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { Grid, Row, Col } from '../components/Grid';
import { ExtFacet } from '../components/ExtFacet';
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
    handleExtToggle = (selectedExt: string) => {
        const { filterParams: { ext } } = this.props;
        Actions.search(this.props.dispatch, this.props.query, {
            ext: ext.find(x => x === selectedExt) ? ext.filter(x => x !== selectedExt) : ext.concat(selectedExt)
        });
    };

    render() {
        const { loading, filterParams, result, facets } = this.props;
        return (
            <Grid>
                <Row>
                    <Col size='is3'>
                        <ExtFacet facet={facets.facets['ext']} selected={filterParams.ext} onToggle={this.handleExtToggle}>
                        </ExtFacet>
                        <FullRefsFacet facets={facets.fullRefsFacet}>
                        </FullRefsFacet>
                    </Col>
                    <Col size='is9'>
                        <Grid>
                            <Row>
                                <Col size='is12'>
                                    {loading ?
                                        <h4><MDSpinner /></h4>
                                        :
                                        <h4>We’ve found {result.size}&nbsp;code results {result.time > 0 ? `(${Math.round(result.time * 1000) / 1000} seconds)` : ''}</h4>
                                    }
                                </Col>
                            </Row>
                            <hr />
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