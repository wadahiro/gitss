import * as React from 'react';

import { Size } from './Modifiers';

const BTag = require('re-bulma/lib/elements/tag').default;

interface TagProps {
    size?: Size;
    style?: any;
    children?: React.ReactElement<any>;
}

const defaultStyle = {
    color: '#666',
    backgroundColor: '#eee'
};

export class Tag extends React.PureComponent<TagProps, void> {
    render() {
        const style = Object.assign({}, defaultStyle, this.props.style);
        return <BTag {...this.props} style={style}>
            {this.props.children}
        </BTag>;
    }
}