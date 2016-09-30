import * as React from 'react';

const style = require('./style.css');
const F = require('react-flexbox-grid');

export function Grid(props) {
    return (
        <F.Grid className={style.grid} {...props}>
            {props.children}
        </F.Grid>
    );
}

export function Row(props) {
    return (
        <F.Row className={style.row} {...props}>
            {props.children}
        </F.Row>
    );
}

export function Col(props) {
    return (
        <F.Col {...props}>
            {props.children}
        </F.Col>
    );
}