import * as React from 'react';

import { Size } from './Modifiers';

const BTag = require('re-bulma/lib/elements/tag').default;

interface TagProps {
    size?: Size;
    style?: any;
    children?: React.ReactElement<any>;
}

export function Tag(props: TagProps) {
    return <BTag {...props}>{props.children}</BTag>
}