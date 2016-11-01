import * as React from 'react';

import { Size } from './Modifiers';

const BTag = require('re-bulma/lib/elements/tag').default;

interface TagProps {
    size?: Size;
    style?: any;
    children?: React.ReactElement<any>;
}

const defaultStyle = {
    color: '#fff',
    backgroundColor: '#3572b0'
};

export function Tag(props: TagProps) {
    const style = Object.assign({}, defaultStyle, props.style);
    return <BTag {...props} style={style}>{props.children}</BTag>
}