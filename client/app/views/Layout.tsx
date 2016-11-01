import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { NavBar } from '../components/NavBar';
import { Container, Grid, Row, Col } from '../components/Grid';
import { RootState, SearchResult } from '../reducers';
import * as Actions from '../actions';

interface Props {
    dispatch?: Dispatch<Action>;
    loading: boolean;
    result: SearchResult;
}

class Layout extends React.Component<Props, void> {

    handleKeyDown = (e: React.KeyboardEvent) => {
        // e.preventDefault();

        Actions.setQuery(this.props.dispatch, e.target['value']);

        if (e.keyCode === 13) {
            Actions.search(this.props.dispatch, e.target['value']);
        }
    };

    render() {
        return (
            <div>
                <NavBar onKeyDown={this.handleKeyDown} loading={this.props.loading} result={this.props.result} />
                <Container style={{ marginTop: 120 }}>
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

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        loading: state.app.present.loading,
        result: state.app.present.result
    };
}

const LayoutContainer = connect(
    mapStateToProps
)(Layout);

export default LayoutContainer;