grammar authz;

// Keywords
VERSION: 'version';
VERSIONNO: [0-9]+;
SERVICE: 'service';
ACCESS: PUBLIC | PRIVATE;
PUBLIC: 'public';
PRIVATE: 'private';
ASSET: 'asset';
PERMISSION: 'permission';
IMPORT: 'import';
DEPENDS: 'depends on';
AS: 'as';
AND: 'and';
OR: 'or';
RESOLVE: '.';
EXPAND: ':';
SEPARATOR: '/';
LBRACE: '{';
RBRACE: '}';

// Tokens
NAME: [a-z][a-z_]*;

WS: [ \r\n\t]+ -> skip;

// Rules
file: statement+;

statement: VERSION VERSIONNO #Version
    | SERVICE NAME #Service
    | IMPORT service=NAME (AS alias=NAME)? #Import
    | permission #GlobalPermission
    | ACCESS? ASSET NAME LBRACE dependency* (permission | computedPermission)+ RBRACE #Asset
;

dependency: DEPENDS service=NAME SEPARATOR asset=NAME AS alias=NAME;
computedPermission: ACCESS? PERMISSION NAME EXPAND permissionExpression;
permission: ACCESS? PERMISSION NAME;

permissionExpression: (internalPermissionRef | externalPermissionRef) #Unary
    | permissionExpression AND permissionExpression #AND
    | permissionExpression OR permissionExpression #OR;

externalPermissionRef: alias=NAME RESOLVE permissionName=NAME;
internalPermissionRef: NAME;