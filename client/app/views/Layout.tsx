import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { NavBar } from '../components/NavBar';
import { Container, Grid, Row, Col } from '../components/Grid';
import { RootState } from '../reducers';
import * as Actions from '../actions';

interface Props {
    dispatch: Dispatch<Action>
}

class Layout extends React.Component<Props, void> {

    handleKeyDown = (e: KeyboardEvent) => {
        // e.preventDefault();

        if (e.keyCode === 13) {
            Actions.search(this.props.dispatch, e.target['value']);
        }
    };

    render() {
        return (
            <div>
                <NavBar onKeyDown={this.handleKeyDown} />
                <Container style={{marginTop: 60}}>
                    <Row>
                        <Col size='is12'>
                            {this.props.children}
                        </Col>
                    </Row>
                </Container>
            </div>
        );
    }
}

function mapStateToProps(state: RootState, props: Props): any {
    return {};
}

const LayoutContainer = connect(
        mapStateToProps
)(Layout);

export default LayoutContainer;