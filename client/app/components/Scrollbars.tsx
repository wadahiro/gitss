import * as React from 'react';

const RCScrollbars = require('react-custom-scrollbars').default;

interface ScrollbarsProps {
    style: any;
    autoHeight?: boolean;
}

export class Scrollbars extends React.PureComponent<ScrollbarsProps, void> {
    render() {
        return (
            <RCScrollbars {...this.props}>{this.props.children}</RCScrollbars>
        );
    }
}