import * as React from 'react';

import { Tag } from './Tag';
import { Panel, PanelHeading, PanelBlock } from './Panel';
import { Menu, MenuLabel, MenuList, MenuLink } from './Menu';
import { Facet, OranizationFacet } from '../reducers';

interface FullRefsFacetProps {
    facets: OranizationFacet[];
}

export function FullRefsFacet(props: FullRefsFacetProps) {
    if (!props.facets || props.facets.length === 0) {
        return null
    }

    return (
        <Panel>
            <PanelHeading>Filter</PanelHeading>
            <Menu style={{ padding: '5px 10px' }}>
                {props.facets.map(organization => {
                    return (
                        <div>
                            <MenuLabel><MenuLink count={organization.count}>{organization.term}</MenuLink></MenuLabel>
                            {organization.projects.length > 0 &&
                                <MenuList>
                                    {organization.projects.map(project => {
                                        return (
                                            <div>
                                                <li><MenuLink count={project.count}>{project.term}</MenuLink></li>
                                                {project.repositories.length > 0 &&
                                                    <li>
                                                        <MenuList>
                                                            {project.repositories.map(repository => {
                                                                return (
                                                                    <div>
                                                                        <li><MenuLink count={repository.count}>{repository.term}</MenuLink></li>
                                                                        {repository.refs.length > 0 &&
                                                                            <li>
                                                                                <MenuList>
                                                                                    {repository.refs.map(ref => {
                                                                                        return (
                                                                                            <li><MenuLink count={ref.count}>{ref.term}</MenuLink></li>
                                                                                        );
                                                                                    })}
                                                                                </MenuList>
                                                                            </li>
                                                                        }
                                                                    </div>
                                                                );
                                                            })}
                                                        </MenuList>
                                                    </li>
                                                }
                                            </div>
                                        );
                                    })}
                                </MenuList>
                            }
                        </div>
                    );
                })}
            </Menu>
        </Panel>
    );
}
