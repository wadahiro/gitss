import * as React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import TextField from 'material-ui/TextField';
import Divider from 'material-ui/Divider';
import RaisedButton from 'material-ui/RaisedButton';

import { Grid, Row, Col } from '../components/Grid';
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
            <div>
                <Row>
                    <Col xs={12}>
                        {loading ?
                            <h4><MDSpinner /></h4>
                            :
                            <h4>We’ve found {result.size}&nbsp;code results {result.time > 0 ? `(${Math.round(result.time * 1000) / 1000} seconds)` : ''}</h4>
                        }
                    </Col>
                </Row>
                <Divider />
                {result.hits.map(x => {
                    return (
                        <Row key={x._source.blob}>
                            <Col xs={12}>
                                <FileContent metadata={x._source.metadata} preview={x.preview} />
                            </Col>
                        </Row>
                    );
                })}
            </div>
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