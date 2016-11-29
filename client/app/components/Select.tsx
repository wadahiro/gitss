import * as React from 'react';
import 'react-select/dist/react-select.css';

const RSelect = require('react-select');

interface SelectProps {
    name?: string;
    className?: string;
    options?: Option[];
    value?: string | string[];
    multi?: boolean;
    onChange?: (values: Option[] | Option) => void;
    icon?: string;
    clearable?: boolean;
    inputProps?: any;
    valueComponent?: any;
    optionComponent?: any;
    arrowRenderer?: any;
    style?: any;
    placeholder?: any;
}

export interface Option {
    label: string;
    value: string;
}

export interface IconOption extends Option {
    icon: string;
}

function isIconOption(o: any): o is IconOption {
    return o && typeof o.icon === 'string';
}

export class Select extends React.PureComponent<SelectProps, void> {
    render() {
        return <RSelect {...this.props} optionComponent={IconOptionComponent} />;
    }
}

export function createIconValue(icon: string, style: Object) {
    return class IconValue extends React.PureComponent<void, void>{
        render() {
            return (
                <div className='Select-value'>
                    <span className='Select-value-label' style={style}>
                        {icon && icon.length > 0 &&
                            <i className={icon} style={{ paddingRight: 5 }} />
                        }
                        {this.props.children}
                    </span>
                </div>
            );
        }
    }
}

interface IconOptionProps {
    title?: string;
    icon?: string;
    className?: string;
    option?: any;
    onSelect?: (option: any, e: any) => void;
    onFocus?: (option: any, e: any) => void;
    isFocused?: boolean;
}

class IconOptionComponent extends React.PureComponent<IconOptionProps, void>{
    handleMouseDown = (event) => {
        event.preventDefault();
        event.stopPropagation();
        this.props.onSelect(this.props.option, event);
    };

    handleMouseEnter = (event) => {
        this.props.onFocus(this.props.option, event);
    };

    handleMouseMove = (event) => {
        if (this.props.isFocused) return;
        this.props.onFocus(this.props.option, event);
    };

    render() {
        let iconStyle = {
            marginRight: 10
        };

        if (isIconOption(this.props.option)) {
            return (
                <div className={this.props.className}
                    onMouseDown={this.handleMouseDown}
                    onMouseEnter={this.handleMouseEnter}
                    onMouseMove={this.handleMouseMove}
                    title={this.props.title}>
                    <i className={this.props.option.icon} style={iconStyle} />
                    {this.props.children}
                </div>
            );

        } else {
            return (
                <div className={this.props.className}
                    onMouseDown={this.handleMouseDown}
                    onMouseEnter={this.handleMouseEnter}
                    onMouseMove={this.handleMouseMove}
                    title={this.props.option.title}>
                    {this.props.children}
                </div>
            );
        }
    }
}
