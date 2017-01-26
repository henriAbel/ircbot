interface ILink {
    Key: number,
    Link: string,
    Sender_id: number,
    Date: number,
    Link_type: string,
    Sender_name: string
}

interface ILoginResponse {
    token: string
}

interface IStatGroupLink {
    Count: number,
    Type: string
}

interface IStatGroupUser {
	Count: number,
	User_name: string,
	Sender_id: number
}

interface IStatDuplicate {
	Count: number,
	User_name: string,
	Sender_id: number
}


interface IStats {
    GroupLink: IStatGroupLink[],
    GroupUser: IStatGroupUser[],
    Duplicates: IStatDuplicate[]
}
