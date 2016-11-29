import * as React from 'react';

import { SizeNum } from './Modifiers';

const BTitle = require('re-bulma/lib/elements/title').default;
const BSubTitle = require('re-bulma/lib/elements/subtitle').default;

interface TitleProps {
    size?: SizeNum;
    style?: any;
    onClick?: (e: any) => void;
}

export class Title extends React.PureComponent<TitleProps, void> {
    render() {
        return <BTitle {...this.props}>
            {this.props.children}
        </BTitle>;
    }
}

export class SubTitle extends React.PureComponent<TitleProps, void> {
    render() {
        return <BSubTitle {...this.props}>
            {this.props.children}
        </BSubTitle>;
    }
}