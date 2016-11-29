import * as React from 'react';
import * as ReactDOM from 'react-dom';

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
    onButtonClick?: (value: string) => void;
}

export class InputTextAddon extends React.PureComponent<InputTextAddonProps, void> {
    input: any = null;

    handleButtonClick = (e: React.SyntheticEvent) => {
        this.props.onButtonClick(ReactDOM.findDOMNode(this.input)['value']);
    };

    render() {
        const { onButtonClick, ...rest } = this.props as InputTextAddonProps;
        return (
            <Addons>
                <InputText ref={(input) => { this.input = input; } } {...rest} />
                <Button size={this.props.size} icon={this.props.buttonIcon} onClick={this.handleButtonClick}>
                    {this.props.buttonTitle}
                </Button>
            </Addons>
        );
    }
}