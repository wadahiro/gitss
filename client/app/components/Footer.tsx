import * as React from 'react';

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
