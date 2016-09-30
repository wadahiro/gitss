import * as React from 'react';
import { Link } from 'react-router'

import { NavBar } from '../components/NavBar';
import { Grid, Row, Col } from '../components/Grid';

export default class Layout extends React.Component<any, void> {

    render() {
        return (
            <div>
                <NavBar />
                <Grid>
                    <Row>
                        <Col xs={12}>
                            {this.props.children}
                        </Col>
                    </Row>
                </Grid>
            </div>
        );
    }
}
