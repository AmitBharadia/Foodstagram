package main

type likecount struct{
	Count int64
}

type like struct{
	User_id string
}

type user_comment struct{
	User_id 	string
	User_name	string
	Timestamp 	int64
	Comment 	string
}

type comment_list struct{
	Comments [] string
}

type data_struct struct{
	Photo_id 		string
	Like_count		int64
	Comment_count	int64
	Likes			[]string
	Comments 	[]user_comment
}

type modal_struct struct{
	Liked		int64
	Comments 	[]user_comment
}

type sns_struct struct{
	Id 		string
	Num 	int64
}
