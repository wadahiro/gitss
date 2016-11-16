import * as React from 'react';

const BFooter = require('re-bulma/lib/layout/footer').default;

interface FooterProps extends React.HTMLAttributes {
}

const defaultStyle = {
    color: '#fff',
    backgroundColor: '#3572b0'
};

export class Footer extends React.PureComponent<FooterProps, void> {
    render() {
        return <BFooter {...this.props} >{this.props.children}</BFooter>;
    }
}
