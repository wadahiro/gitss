import * as React from 'react';
import 'react-select/dist/react-select.css';

const RSelect = require('react-select');

interface SelectProps {
    name?: string;
    className?: string;
    options?: Option[];
    value?: string | string[];
    multi?: boolean;
    onChange?: (values: Option[]) => void;
    icon?: string;
    clearable?: boolean;
    valueComponent?: any;
    arrowRenderer?: any;
    style?: any;
    placeholder?: any;
}

export interface Option {
    label: string;
    value: string;
}

export class Select extends React.PureComponent<SelectProps, void> {
    render() {
        return <RSelect {...this.props} />;
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

class IconOption extends React.PureComponent<IconOptionProps, void>{
    handleMouseDown(event) {
        event.preventDefault();
        event.stopPropagation();
        this.props.onSelect(this.props.option, event);
    }
    handleMouseEnter(event) {
        this.props.onFocus(this.props.option, event);
    }
    handleMouseMove(event) {
        if (this.props.isFocused) return;
        this.props.onFocus(this.props.option, event);
    }
    render() {
        let gravatarStyle = {
            borderRadius: 3,
            display: 'inline-block',
            marginRight: 10,
            position: 'relative',
            top: -2,
            verticalAlign: 'middle',
        };
        return (
            <div className={this.props.className}
                onMouseDown={this.handleMouseDown}
                onMouseEnter={this.handleMouseEnter}
                onMouseMove={this.handleMouseMove}
                title={this.props.option.title}>
                <i className='fa fa-th-large' />
                {this.props.children}
            </div>
        );
    }
}