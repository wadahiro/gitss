import * as React from 'react';

const BContainer = require('re-bulma/lib/layout/container').default;
const Columns = require('re-bulma/lib/grid/columns').default;
const Column = require('re-bulma/lib/grid/column').default;
const BSection = require('re-bulma/lib/layout/section').default;
const BContent = require('re-bulma/lib/elements/content').default;

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

export function Section(props) {
    return (
        <BSection {...props}>
            {props.children}
        </BSection>
    );
}

interface RowProps {
    isGapless?: boolean;
    children?: React.ReactElement<any>;
    style?: Object;
}

export function Row(props: RowProps) {
    return (
        <Columns {...props}>
            {props.children}
        </Columns>
    );
}

interface ColProps {
    size?: 'is1' | 'is2' | 'is3' | 'is4' | 'is5' | 'is6' | 'is7' | 'is8' | 'is9' | 'is10' | 'is11' | 'is12' | 'isNarrow';
    children?: React.ReactElement<any>;
    style?: any;
}

export function Col(props: ColProps) {
    return (
        <Column {...props}>
            {props.children}
        </Column>
    );
}

const tableRowStyle = {
    wrapper: {
        display: 'table'
    },
    item: {
        display: 'table-cell',
        textAlign: 'left'
    }
};

interface TRowProps {
    style?: Object;
    children?: any;
}

export function TRow(props: TRowProps) {
    let style = tableRowStyle.wrapper;
    if (props && props.style) {
        style = {
            ...style,
            ...props.style
        };
    }
    return (
        <div style={style}>
            {props.children}
        </div>
    );
}

interface TColProps {
    style?: Object;
    children?: any;
}

export function TCol(props: TColProps) {
    let style = tableRowStyle.item;
    if (props && props.style) {
        style = {
            ...style,
            ...props.style
        };
    }
    return (
        <div style={style}>
            {props.children}
        </div>
    );
}

export function Content(props) {
    return (
        <BContent {...props} />
    );
}

interface StickyFooterPageProps {
    footer: React.ReactElement<any>;
}

export class StickyFooterPage extends React.PureComponent<StickyFooterPageProps, void> {
    render() {
        return (
            <div style={stickyFooterPageStyle}>
                <div style={stickyFooterPageMainStyle}>
                    {this.props.children}
                </div>
                <div>
                    {this.props.footer}
                </div>
            </div>
        );
    }
}

const stickyFooterPageStyle = {
    minHeight: '100vh',
    display: 'flex',
    flexDirection: 'column'
};

const stickyFooterPageMainStyle = {
    flex: 1
};