import * as React from 'react';

const BContainer = require('re-bulma/lib/layout/container').default;
const Columns = require('re-bulma/lib/grid/columns').default;
const Column = require('re-bulma/lib/grid/column').default;

export function Container(props) {
    return (
        <BContainer {...props}>
            {props.children}
        </BContainer>
    );
}

export function Grid(props) {
    return (
        <div {...props}>
            {props.children}
        </div>
    );
}

export function Row(props) {
    return (
        <Columns {...props}>
            {props.children}
        </Columns>
    );
}

interface ColProps {
    size: 'is1' | 'is2' | 'is3' | 'is4' | 'is5' | 'is6' | 'is7' | 'is8' | 'is9' | 'is10' | 'is11' | 'is12';
    children?: React.ReactElement<any>;
}

export function Col(props: ColProps) {
    return (
        <Column {...props}>
            {props.children}
        </Column>
    );
}