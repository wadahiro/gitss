import * as React from 'react';

import { Size } from './Modifiers';

const BPanel = require('re-bulma/lib/components/panel/panel').default;
const BPanelHeading = require('re-bulma/lib/components/panel/panel-heading').default;
const BPanelBlock = require('re-bulma/lib/components/panel/panel-block').default;

interface PanelProps extends React.HTMLAttributes {
}

export class Panel extends React.PureComponent<PanelProps, void> {
    render() {
        return <BPanel {...this.props}>{this.props.children}</BPanel>;
    }
}

export class PanelHeading extends React.PureComponent<PanelProps, void> {
    render() {
        const style = Object.assign({},
            {
                fontWeight: 600,
                padding: '5px 10px'
            },
            this.props.style);
        return <BPanelHeading {...this.props} style={style}>{this.props.children}</BPanelHeading>;
    }
}

interface PanelBlockProps extends React.HTMLAttributes {
    isActive?: boolean;
    icon?: string;
}

export class PanelBlock extends React.PureComponent<PanelBlockProps, void> {
    render() {
        return <BPanelBlock {...this.props}>{this.props.children}</BPanelBlock>;
    }
}