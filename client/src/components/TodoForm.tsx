import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const TodoForm = () => {
    const [newTodo, setNewTodo] = useState("");
    const queryClient = useQueryClient();

    const { mutate: createTodo, isPending: isCreating } = useMutation({
        mutationKey: ["createTodo"],
        mutationFn: async (todoTitle: string) => {  // 改名為 todoTitle 更明確
            const res = await fetch('http://localhost:3000/todos', {  // 確認後端 URL 正確
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ 
                    title: todoTitle,      // 改這裡：title 而非 body
                    completed: false       // 後端預設值
                }),
            });
            const data = await res.json();
            if (!res.ok) {
                throw new Error(data.error || "Something went wrong");
            }
            return data;
        },
        onSuccess: () => {
            setNewTodo("");
            queryClient.invalidateQueries({ queryKey: ["todos"] });
        },
        onError: (error: any) => {
            alert(error.message);
        },
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (newTodo.trim()) {  // 避免空字串
            createTodo(newTodo);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                value={newTodo}
                onChange={(e) => setNewTodo(e.target.value)}
                placeholder="新增 Todo"
                disabled={isCreating}
            />
            <button type="submit" disabled={isCreating || !newTodo.trim()}>
                {isCreating ? "新增中..." : "新增"}
            </button>
        </form>
    );
};
export default TodoForm;
