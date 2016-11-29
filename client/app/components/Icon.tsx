import * as React from 'react';

import { Size } from './Modifiers';

const BIcon = require('re-bulma/lib/elements/icon').default;

interface IconProps extends React.DOMAttributes {
    icon: string;
    onClick?: () => void;
}

export class Icon extends React.PureComponent<IconProps, void> {
    render() {
        const style = this.props.onClick ? { cursor: 'pointer' } : {};
        return <i style={style} className={`fa fa-${this.props.icon}`} onClick={this.props.onClick} />;
    }
}

interface IconLinkProps extends IconProps {
    href?: string;
    onClick?: () => void;
}

export class IconLink extends React.PureComponent<IconLinkProps, void> {
    static defaultProps = {
        href: '',
        onClick: () => { }
    };

    handleClick = (e: React.SyntheticEvent) => {
        e.preventDefault();
        this.props.onClick();
    };

    render() {
        return <a href={this.props.href} onClick={this.props.onClick}><Icon icon={this.props.icon} /></a>;
    }
}