from typing import List, Dict, Any, Optional


class Command:
    """
    Represents a command specification for carapace-sh/carapace-spec.

    Attributes:
        name (str): The name of the command.
        description (str): The description of the command.
        flags (Dict[str, Dict[str, Any]]): The flags for the command.
        mark_flags_one_required (Optional[List[str]]): List of flag names where at least one is required.
    """

    def __init__(
        self,
        name: str,
        description: str = "",
        flags: Optional[Dict[str, Dict[str, Any]]] = None,
        mark_flags_one_required: Optional[List[str]] = None,
    ) -> None:
        self.name = name
        self.description = description
        self.flags = flags if flags is not None else {}
        self.mark_flags_one_required = mark_flags_one_required or []

    def add_flag(self, flag_name: str, flag_spec: Dict[str, Any]) -> None:
        """
        Adds a flag to the command.

        Args:
            flag_name (str): The name of the flag.
            flag_spec (Dict[str, Any]): The specification for the flag.
        """
        self.flags[flag_name] = flag_spec

    def mark_flags_one_required_set(self, flag_names: List[str]) -> None:
        """
        Sets the flags for which at least one is required.

        Args:
            flag_names (List[str]): List of flag names.
        """
        self.mark_flags_one_required = flag_names

    def to_spec(self) -> Dict[str, Any]:
        """
        Serializes the command to a spec dictionary.

        Returns:
            Dict[str, Any]: The serialized command spec.
        """
        spec: Dict[str, Any] = {
            "name": self.name,
            "description": self.description,
            "flags": self.flags,
        }
        if self.mark_flags_one_required:
            spec["markFlagsOneRequired"] = self.mark_flags_one_required
        return spec


# Example usage and test
def test_command_mark_flags_one_required() -> None:
    """
    Test that Command correctly serializes markFlagsOneRequired.
    """
    cmd = Command(
        name="example",
        description="An example command",
        flags={
            "--foo": {"type": "string", "description": "foo flag"},
            "--bar": {"type": "bool", "description": "bar flag"},
        },
    )
    cmd.mark_flags_one_required_set(["--foo", "--bar"])
    spec = cmd.to_spec()
    assert "markFlagsOneRequired" in spec
    assert spec["markFlagsOneRequired"] == ["--foo", "--bar"]
    assert spec["flags"]["--foo"]["type"] == "string"
    assert spec["flags"]["--bar"]["type"] == "bool"


if __name__ == "__main__":
    test_command_mark_flags_one_required()