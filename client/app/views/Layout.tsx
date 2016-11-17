import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { Footer } from '../components/Footer';
import { Container, Grid, Row, Col } from '../components/Grid';
import { RootState, SearchResult } from '../reducers';
import * as Actions from '../actions';

export default class Layout extends React.PureComponent<void, void> {
    render() {
        return (
            <div>
                {this.props.children}
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
