import * as React from 'react';

import { SizeNum } from './Modifiers';

const BTitle = require('re-bulma/lib/elements/title').default;

interface TitleProps {
    size?: SizeNum;
    style?: any;
    children?: React.ReactElement<any>;
    onClick?: (e: any) => void;
}

export class Title extends React.PureComponent<TitleProps, void> {
    render() {
        return <BTitle {...this.props}>
            {this.props.children}
        </BTitle>;
    }
}