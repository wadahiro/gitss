import * as React from 'react';
import 'react-select/dist/react-select.css';

const RSelect = require('react-select');

interface SelectProps {
    name?: string;
    options?: Option[];
    value?: string[];
    multi?: boolean;
    onChange?: (values: Option[]) => void;
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