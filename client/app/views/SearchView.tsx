import * as React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { Grid, Row, Col } from '../components/Grid';
import { FullRefsFacet } from '../components/FullRefsFacet';
import { Facets } from '../components/Facets';
import { FileContent } from '../components/FileContent';
import { RootState, SearchResult } from '../reducers';

const MDSpinner = require('react-md-spinner').default;

interface Props {
    loading: boolean;
    result: SearchResult;
}

class SearchView extends React.Component<Props, void> {
    render() {
        const { loading, result } = this.props;
        return (
            <Grid>
                <Row>
                    <Col size='is3'>
                        <FullRefsFacet facets={result.fullRefsFacet}>
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
                                            <FileContent metadata={x._source.metadata} preview={x.preview} />
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
        result: state.app.present.result
    };
}

const SearchViewContainer = connect(
    mapStateToProps
)(SearchView);

export default SearchViewContainer;