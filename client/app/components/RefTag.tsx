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

export class RefTag extends React.PureComponent<RefTagProps, void> {
    render() {
        const style = Object.assign({}, defaultStyle, this.props.style);
        return <Tag {...this.props} style={style}>
            <Icon {...this.props} style={{ marginRight: 3 }} />
            {this.props.children}
        </Tag>;
    }
}

export class Icon extends React.PureComponent<RefTagProps, void> {
    render() {
        if (this.props.type === 'branch') {
            return <i className='fa fa-code-fork' style={this.props.style} />;
        } else {
            return <i className='fa fa-tag' style={this.props.style} />;
        }
    }
}