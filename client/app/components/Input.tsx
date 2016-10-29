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

export function InputText(props: InputTextProps) {
    return <BInput {...props} type='text'>{props.children}</BInput>
}