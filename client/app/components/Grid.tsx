import * as React from 'react';

// const style = require('./style.css');
// const F = require('reflexbox');

const style = {
    grid: "grid",
    row: "row"
}

export function Grid(props) {
    return (
        <div className='container' {...props}>
            {props.children}
        </div>
    );
}

export function Row(props) {
    return (
        <div className='row' {...props}>
            {props.children}
        </div>
    );
}

export function Col(props) {
    return (
        <div className={`col-${props.xs}`}>
            {props.children}
        </div>
    );
}