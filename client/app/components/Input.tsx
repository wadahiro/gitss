import * as React from 'react';

import { Size } from './Modifiers';
import { Button } from './Button';
import { Addons } from './Addons';

const BInput = require('re-bulma/lib/forms/input').default;

interface InputTextProps extends React.DOMAttributes {
    placeholder?: string;
    icon?: string;
    size?: Size;
    hasIcon?: boolean;
    style?: any;
    value?: string;
    defaultValue?: string;
    isExpanded?: boolean;
}

export class InputText extends React.PureComponent<InputTextProps, void> {
    render() {
        return <BInput {...this.props} type='text'>{this.props.children}</BInput>;
    }
}

interface InputTextAddonProps extends InputTextProps {
    buttonTitle?: string;
    buttonIcon?: string;
}

export class InputTextAddon extends React.PureComponent<InputTextAddonProps, void> {
    render() {
        return (
            <Addons>
                <InputText {...this.props} />
                <Button size={this.props.size} icon={this.props.buttonIcon}>
                    {this.props.buttonTitle}
                </Button>
            </Addons>
        );
    }
}