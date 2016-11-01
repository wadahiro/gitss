import * as React from 'react';

const BPagination = require('re-bulma/lib/components/pagination/pagination').default;
const BPageButton = require('re-bulma/lib/components/pagination/page-button').default;

interface PaginationProps extends React.HTMLAttributes {
    children?: React.ReactElement<any>;
    isActive?: boolean;
}

const defaultStyle = {
    color: '#fff',
    backgroundColor: '#3572b0'
};

export function Pagination(props: PaginationProps) {
    return <BPagination {...props} >{props.children}</BPagination>;
}

export function PageButton(props: PaginationProps) {
    let style = props.style;
    if (props.isActive) {
        style = Object.assign({}, defaultStyle, props.style);
    }

    return <BPageButton {...props} style={style}>{props.children}</BPageButton>;
}

