import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { NavBar } from '../components/NavBar';
import { Footer } from '../components/Footer';
import { Container, Grid, Row, Col } from '../components/Grid';
import { RootState, SearchResult } from '../reducers';
import * as Actions from '../actions';

const BHero = require('re-bulma/lib/layout/hero').default;
const BHeroHead = require('re-bulma/lib/layout/hero-head').default;
const BHeroBody = require('re-bulma/lib/layout/hero-body').default;

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
        const { result } = this.props;

        return (
            <div>
                <NavBar onKeyDown={this.handleKeyDown} loading={this.props.loading} result={this.props.result} />
                <Container style={{ marginTop: 80 }}>
                    <Row>
                        <Col size='is12'>
                            {this.props.children}
                        </Col>
                    </Row>
                </Container>
                <Footer>
                    <Container>
                        <p style={{ textAlign: 'center' }}>
                            <strong>GitSS</strong> - Git Source Search. The source code is licensed&nbsp;
        <a href="http://opensource.org/licenses/mit-license.php">MIT</a>.
      </p>
                        <p style={{ textAlign: 'center' }}>
                            <a className="icon" href="https://github.com/wadahiro/gitss">
                                <i className="fa fa-github"></i>
                            </a>
                        </p>
                    </Container>
                </Footer>
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