import * as React from 'react';

import { Size } from './Modifiers';

const BInput = require('re-bulma/lib/forms/input').default;

interface InputTextProps extends React.DOMAttributes {
    placeholder?: string;
    icon?: string;
    size?: Size;
    hasIcon?: boolean;
    style?: any;
}

export class InputText extends React.PureComponent<InputTextProps, void> {
    render() {
        return <BInput {...this.props} type='text'>{this.props.children}</BInput>;
    }
}