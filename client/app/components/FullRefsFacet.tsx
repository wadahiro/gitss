import * as React from 'react';

import { Tag } from './Tag';
import { Facet } from '../reducers';

const BMenu = require('re-bulma/lib/components/menu/menu').default;
const BMenuLabel = require('re-bulma/lib/components/menu/menu-label').default;
const BMenuList = require('re-bulma/lib/components/menu/menu-list').default;
const BMenuLink = require('re-bulma/lib/components/menu/menu-link').default;

interface FullRefsFacetProps {
    facet: Facet;
}

function isRef(term: string) {
    return term.includes(':') && term.includes('/') && term.indexOf(':') !== term.lastIndexOf(':');
}

function isRepository(term: string) {
    return term.includes(':') && term.includes('/') && term.indexOf(':') === term.lastIndexOf(':');
}

function isProject(term: string) {
    return term.includes(':') && !term.includes('/');
}

function isOrganization(term: string) {
    return !term.includes(':');
}


export function FullRefsFacet(props: FullRefsFacetProps) {
    if (!props.facet) {
        return null;
    }

    let prev = '';
    let organizations = []
    let projects = [];
    let repositories = [];
    let refs = [];

    const terms = props.facet.terms.reverse();

    terms.forEach(term => {
        if (isRef(prev) && isRepository(term.term)) {
            // repository name           
            repositories.push(
                <li key={term.term}><MenuLink count={term.count}>{term.term.split('/')[1]}</MenuLink></li>
            );
            // refs
            repositories.push(
                <li>
                    <BMenuList>{refs}</BMenuList>
                </li>
            );
            refs = [];
        }
        if (isRepository(prev) && isProject(term.term)) {
            // project name             
            projects.push(
                <li key={term.term}><MenuLink count={term.count}>{term.term.split(':')[1]}</MenuLink></li>
            );
            // repositories
            projects.push(
                <li>
                    <BMenuList>{repositories}</BMenuList>
                </li>
            );
            repositories = [];
        }
        if (isProject(prev) && isOrganization(term.term)) {
            // organization name
            organizations.push(
                <BMenuLabel>{term.term}</BMenuLabel>
            );
            // projects
            organizations.push(
                <BMenuList>{projects}</BMenuList>
            );
            projects = [];
        }

        // ref
        if (isRef(term.term)) {
            refs.push(
                <li key={term.term}><MenuLink count={term.count}>{term.term.split(':')[2]}</MenuLink></li>
            );
        }

        prev = term.term;
    })


    return (
        <BMenu {...props}>
            {organizations}
        </BMenu>
    );
}

interface MenuLinkProps {
    count: number;
    isActive?: boolean;
    children?: React.ReactElement<any>
}

export function MenuLink(props: MenuLinkProps) {
    return <BMenuLink {...props}>{props.children}<Tag size='isSmall' style={{ float: 'right' }}>{props.count}</Tag></BMenuLink>
}