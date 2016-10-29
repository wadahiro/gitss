import * as React from 'react';

import { Tag } from './Tag';
import { Facet, OranizationFacet } from '../reducers';

const BMenu = require('re-bulma/lib/components/menu/menu').default;
const BMenuLabel = require('re-bulma/lib/components/menu/menu-label').default;
const BMenuList = require('re-bulma/lib/components/menu/menu-list').default;
const BMenuLink = require('re-bulma/lib/components/menu/menu-link').default;

interface FullRefsFacetProps {
    facets: OranizationFacet[];
}

export function FullRefsFacet(props: FullRefsFacetProps) {
    if (!props.facets || props.facets.length === 0) {
        return null
    }

    return (
        <BMenu>
            {props.facets.map(organization => {
                return (
                    <div>
                        <BMenuLabel><MenuLink count={organization.count}>{organization.term}</MenuLink></BMenuLabel>
                        {organization.projects.length > 0 &&
                            <BMenuList>
                                {organization.projects.map(project => {
                                    return (
                                        <div>
                                            <li><MenuLink count={project.count}>{project.term}</MenuLink></li>
                                            {project.repositories.length > 0 &&
                                                <li>
                                                    <BMenuList>
                                                        {project.repositories.map(repository => {
                                                            return (
                                                                <div>
                                                                    <li><MenuLink count={repository.count}>{repository.term}</MenuLink></li>
                                                                    {repository.refs.length > 0 &&
                                                                        <li>
                                                                            <BMenuList>
                                                                                {repository.refs.map(ref => {
                                                                                    return (
                                                                                        <li><MenuLink count={ref.count}>{ref.term}</MenuLink></li>
                                                                                    );
                                                                                })}
                                                                            </BMenuList>
                                                                        </li>
                                                                    }
                                                                </div>
                                                            );
                                                        })}
                                                    </BMenuList>
                                                </li>
                                            }
                                        </div>
                                    );
                                })}
                            </BMenuList>
                        }
                    </div>
                );
            })}
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