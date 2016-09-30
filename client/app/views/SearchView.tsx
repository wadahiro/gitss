import * as React from 'react';
import { Link } from 'react-router'
import TextField from 'material-ui/TextField';
import RaisedButton from 'material-ui/RaisedButton';

import { Grid, Row, Col } from '../components/Grid';
import { FileContent } from '../components/FileContent';

export default class SearchView extends React.Component<any, void> {
    render() {
        return (
            <Grid>
                <Row>
                    <Col xs={12}>
                        <FileContent />
                    </Col>
                </Row>
            </Grid>
        );
    }
}
