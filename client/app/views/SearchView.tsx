import * as React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import TextField from 'material-ui/TextField';
import Divider from 'material-ui/Divider';
import RaisedButton from 'material-ui/RaisedButton';

import { Grid, Row, Col } from '../components/Grid';
import { FileContent } from '../components/FileContent';
import { RootState, SearchResult } from '../reducers';

interface Props {
    result: SearchResult[]
}

class SearchView extends React.Component<Props, void> {
    render() {
        const { result } = this.props;
        return (
            <Grid>
                <Row>
                    <Col xs={12}>
                        <h2>We’ve found {result.length} code results</h2>
                    </Col>
                </Row>
                <Divider/>
                {result.map(x => {
                    return (
                        <Row key={x.blob}>
                            <Col xs={12}>
                                <FileContent metadata={x.metadata} content={x.content} />
                            </Col>
                        </Row>
                    );
                })}
            </Grid>
        );
    }
}

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        result: state.app.present.result
    };
}

const SearchViewContainer = connect(
    mapStateToProps
)(SearchView);

export default SearchViewContainer;