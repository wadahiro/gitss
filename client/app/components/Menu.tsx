import * as React from 'react';

import { Tag } from './Tag';
import { Size } from './Modifiers';


const BMenu = require('re-bulma/lib/components/menu/menu').default;
const BMenuLabel = require('re-bulma/lib/components/menu/menu-label').default;
const BMenuList = require('re-bulma/lib/components/menu/menu-list').default;
const BMenuLink = require('re-bulma/lib/components/menu/menu-link').default;

interface MenuProps extends React.HTMLAttributes {
}

export function Menu(props: MenuProps) {
    return <BMenu {...props}>{props.children}</BMenu>;
}

export function MenuLabel(props: MenuProps) {
    return <BMenuLabel {...props}>{props.children}</BMenuLabel>;
}

export function MenuList(props: MenuProps) {
    return <BMenuList {...props}>{props.children}</BMenuList>;
}

interface MenuLinkProps extends React.HTMLAttributes {
    count?: number;
    isActive?: boolean;
    isToggled?: boolean;
}

export function MenuLink(props: MenuLinkProps) {
    let toggledStyle = {} as any;
    if (props.isToggled) {
        toggledStyle = {
            borderLeft: '3px solid #d26911',
            paddingLeft: 7,
            backgroundColor: '#f5f7fa'
        };
    }
    const linkStyle = Object.assign({}, toggledStyle, props.style);

    if (typeof props.count === 'number') {
        return <BMenuLink {...props} style={linkStyle}>
            <span style={{ paddingRight: 30 }}>{props.children}</span>
            <Tag size='isSmall' style={{ float: 'right' }}>{props.count}</Tag>
        </BMenuLink>;
    }
    return <BMenuLink {...props} style={linkStyle}>{props.children}</BMenuLink>;
}