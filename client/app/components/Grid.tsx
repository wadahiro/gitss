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

export function Col(props) {
    return (
        <Column>
            {props.children}
        </Column>
    );
}