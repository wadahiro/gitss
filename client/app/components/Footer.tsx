import * as React from 'react';

import { Container, Grid, Row, Col } from '../components/Grid';

const BFooter = require('re-bulma/lib/layout/footer').default;

interface FooterProps extends React.HTMLAttributes {
}

const defaultStyle = {
};

export class Footer extends React.PureComponent<FooterProps, void> {
    render() {
        return <BFooter style={defaultStyle} {...this.props} >{this.props.children}</BFooter>;
    }
}

export class AppFooter extends React.PureComponent<FooterProps, void> {
    render() {
        return (
            <Footer style={footerStyle}>
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
        );
    }
}


const footerStyle = {
    width: '100%',
    borderTop: '1px solid rgb(204, 204, 204)'
    // height: '100px',
    // position: 'absolute',
    // bottom: '0px'
};
