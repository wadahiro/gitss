import * as React from 'react';

import { Size } from './Modifiers';
import { Tag } from './Tag';

interface RefTagProps {
    size?: Size;
    style?: any;
    children?: React.ReactElement<any>;
    type: 'branch' | 'tag';
}

const defaultStyle = {
    color: '#5c5c5c',
    backgroundColor: '#f8fafc',
    border: '1px solid #d3d6db',
    borderRadius: 5,
    paddingLeft: 5,
    paddingRight: 5,
    fontSize: 10,
    margin: 2
};

export function RefTag(props: RefTagProps) {
    const style = Object.assign({}, defaultStyle, props.style);
    return <Tag {...props} style={style}>
        <Icon {...props} style={{ marginRight: 3 }} />
        {props.children}
    </Tag>;
}

function Icon(props) {
    if (props.type === 'branch') {
        return <i className='fa fa-code-fork' style={props.style} />;
    } else {
        return <i className='fa fa-tag' style={props.style} />;
    }
}