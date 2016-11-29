import * as React from 'react';

import { Size, Color } from './Modifiers';

const BButton = require('re-bulma/lib/elements/button').default;

interface ButtonProps {
    color?: Color;
    size?: Size;
    icon?: string;
}

export class Button extends React.PureComponent<ButtonProps, void> {
    render() {
        return <BButton {...this.props} />;
    }
}
