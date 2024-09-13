library IEEE;
use IEEE.STD_LOGIC_1164.all;
use IEEE.NUMERIC_STD.all;

entity NBitAdder_tb is
end NBitAdder_tb;

architecture Behavioral of NBitAdder_tb is

    constant N : integer := 4;

    signal A        : std_logic_vector(N-1 downto 0) := (others => '0');
    signal B        : std_logic_vector(N-1 downto 0) := (others => '0');
    signal Sum      : std_logic_vector(N-1 downto 0);
    signal CarryOut : std_logic_vector(N-1 downto 0);

begin

    uut : entity work.NBitAdder
        generic map (N => N)
        port map (
            A        => A,
            B        => B,
            Sum      => Sum,
            CarryOut => CarryOut
            );

    -- Stimulus process to apply test vectors
    stim_proc : process
    begin
        -- Test Case 1: 0 + 1
        A <= "0000";
        B <= "0001";
        wait for 10 ns;

        -- Test Case 2: 3 + 3
        A <= "0011";
        B <= "0011";
        wait for 10 ns;

        -- Test Case 3: 15 + 1
        A <= "1111";
        B <= "0001";
        wait for 10 ns;

        -- Test Case 4: 10 + 5
        A <= "1010";
        B <= "0101";
        wait for 10 ns;

        -- Test Case 5: Maximum values
        A <= (others => '1');
        B <= (others => '1');
        wait for 10 ns;

        -- Finish simulation
        wait;
    end process;

end Behavioral;
