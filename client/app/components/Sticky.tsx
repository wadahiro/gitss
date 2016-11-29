import * as React from 'react';
import * as RS from 'react-sticky';

interface StickyContainerProps extends React.DOMAttributes {
}

export class StickyContainer extends React.PureComponent<StickyContainerProps, void> {
    render() {
        return <RS.StickyContainer {...this.props} />;
    }
}

interface StickyProps {
    stickyStyle?: Object;
    stickyClassName?: string;
    topOffset?: number;
    bottomOffset?: number;
    className?: string;
    style?: Object;
    onStickyStateChange?: () => void;
    isActive?: boolean;
}

export class Sticky extends React.PureComponent<StickyProps, void> {
    render() {
        return <RS.Sticky {...this.props} />;
    }
}