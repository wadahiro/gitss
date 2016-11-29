import * as React from 'react';

import { Size } from './Modifiers';

const BIcon = require('re-bulma/lib/elements/icon').default;

interface IconProps {
    icon: string;
    onClick?: () => void;
    size?: Size;
}

export class Icon extends React.PureComponent<IconProps, void> {
    render() {
        const style = this.props.onClick ? { cursor: 'pointer' } : {};
        return <BIcon {...this.props} style={style} className={`fa fa-${this.props.icon}`} />;
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