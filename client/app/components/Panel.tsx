import * as React from 'react';

import { Size } from './Modifiers';

const BPanel = require('re-bulma/lib/components/panel/panel').default;
const BPanelHeading = require('re-bulma/lib/components/panel/panel-heading').default;
const BPanelBlock = require('re-bulma/lib/components/panel/panel-block').default;

interface PanelProps extends React.HTMLAttributes {
}

export function Panel(props: PanelProps) {
    return <BPanel {...props}>{props.children}</BPanel>;
}

export function PanelHeading(props: PanelProps) {
    const style = Object.assign({},
        {
            fontWeight: 600,
            padding: '5px 10px'
        },
        props.style);
    return <BPanelHeading {...props} style={style}>{props.children}</BPanelHeading>;
}

interface PanelBlockProps extends React.HTMLAttributes {
    isActive?: boolean;
    icon?: string;
}

export function PanelBlock(props: PanelBlockProps) {
    return <BPanelBlock {...props}>{props.children}</BPanelBlock>;
}